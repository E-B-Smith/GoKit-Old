

drop database if exists AWSBackup;
drop schema if exists AWSBackup cascade;
drop user if exists AWSBackup;


create user AWSBackup 
	with createdb login password 'AWSBackup';


create database AWSBackup
	with encoding 'UTF8' owner AWSBackup;


create schema AWSBackup authorization AWSBackup;


set search_path TO AWSBackup, public;


create table AWSParameterTable
	(
	 version 			varchar(10)
	,localBundlePath	varchar(256)
	,AWSBucketName		varchar(256)
	,AWSAccessKeyID		varchar(64)
	,AWSAccessSecret	varchar(64)
	,AWSRegion			varchar(16)
	);

insert into AWSParameterTable (version) values ('1.00.001');


create type AWSStorageState as enum
	(
	 'AWSStorageLocal'
	,'AWSStorageGlacier'
	,'AWSStorageRestoring'
	,'AWSStorageRestored'
	,'AWSStorageStandard'
	);


create table AWSObjectTable
	(
	 key				varchar(256) unique not null primary key

	,awsDate			timestamptz
	,awsBytes			integer
	,awsStorage			AWSStorageState

	,localDate			timestamptz
	,localBytes			integer
	);


create type AWSLogLevel as enum
	(
	 'AWSLogDebug'
	,'AWSLogInfo'
	,'AWSLogStart'
	,'AWSLogExit'
	,'AWSLogWarning'
	,'AWSLogError'
	);


create table AWSLogTable
	(
	 entry			serial			unique not null primary key
	,time			timestamptz		not null
	,processname	varchar(16)		not null
	,filename		varchar(32)		not null
	,linenumber		integer			not null
	,pid 			integer			not null
	,level			AWSLogLevel		not null
	,message		varchar(512) 	not null
	);
create index AWSLogTimeIndex on AWSLogTable(time, entry);


create table AWSStatusTable
	(
	 entry				serial			unique not null primary key
	,time				timestamptz		not null
	,pid				integer			not null
	,bundle				varchar(255)	
	,updatedBytes		bigint
	,totalUpdateBytes	bigint
	,deletedBytes		bigint
	,totalDeleteBytes	bigint
	,totalLocalBytes	bigint
	,message			varchar(512)
	);
create index AWSStatusPIDIndex on AWSStatusTable(pid);


--	Load local data:


create function AWSBulkLoadLocalData(filepath text) returns integer as
	$$
	begin

	truncate table AWSObjectTable;
	truncate table BulkLoadTable cascade;

	-- create temporary table BulkLoadTable
	create table if not exists BulkLoadTable
		(
		 key				varchar(256) unique not null primary key
		,localDate			timestamptz
		,localBytes			integer
		)
	; -- on commit drop;

	copy BulkLoadTable from 'filepath' with null '';

	lock table AWSObjectTable in exclusive mode;

	update AWSObjectTable
		set localDate = BulkLoadTable.localDate,
			localBytes = BulkLoadTable.localBytes
		from BulkLoadTable
		where BulkLoadTable.key = AWSObjectTable.key;
		
	insert into AWSObjectTable
			(key, localDate, localBytes)
		select
			BulkLoadTable.key,
			BulkLoadTable.localDate,
			BulkLoadTable.localBytes
		from BulkLoadTable
		left outer join AWSObjectTable on 
			(AWSObjectTable.key = BulkLoadTable.key)
		where AWSObjectTable.key is null;
		
	end; 
	$$ 
	language plpgsql
	returns null on null input;


-- Load AWS Data:


create function awsxmlns() returns text[][] as
	$$
	begin
	return array[array['ns', 'http://s3.amazonaws.com/doc/2006-03-01/']];
	end;
	$$ 
	language plpgsql immutable;


create table AWSBulkLoadTable
	(
	loadID		   serial primary key,
	bucket         varchar(255) not null,
	prefix		   varchar(512) not null,
	isTruncated    bool         not null,
	
	constraint AWSBulkLoadTableUniqueConstraint
		unique (bucket, prefix)
	);


create table AWSBulkLoadDataTable
	(
	loadID			integer not null,
	path 			varchar(512) not null,
	awsDate 		timestamptz not null,
	awsBytes 		integer not null,
	awsStorage		varchar(16) not null,
	
	constraint AWSBulkLoadDataTableKeyConstraint
		foreign key (loadID)
		references AWSBulkLoadTable(loadID)
		on delete cascade,
		
	constraint AWSBulkLoadDataTableUniqueConstraint
		unique (loadID, path)
	);
create index AWSBuildLoadDataPathIndex on AWSBulkLoadDataTable(loadID, path);


create function AWSBulkLoadAWSData(filepath text) returns integer as
	$$
	begin

	-- For debugging:
	
	truncate table bulkLoadText cascade;
	truncate table bulkLoadXML cascade;
	truncate table AWSBulkLoadTable cascade;
	truncate table AWSBulkLoadDataTable cascade;

	-- Create a temp text table to copy the raw data into:

	-- create temp table bulkLoadText (textData text) on commit drop;
	create table if not exists  bulkLoadText (textData text);

	copy bulkLoadText (textData)
		from '/Users/Edward/Development/go/src/violent.blue/go/aws-bu/TestData/TestBackup.adata';
	--	from stdin;

	delete from bulkLoadText where ctid in 
		(select ctid from bulkLoadText limit 1);

	-- Create a temp table to hold the parsed XML so we only parse it once:

	-- create temp table bulkLoadXML
	create table if not exists bulkLoadXML
		(
		xmlData  	xml,
		bucket	 	varchar(256),
		prefix	 	varchar(512),
		isTruncated bool
		) 
	; --	on commit drop;

	insert into bulkLoadXML
		with parsedXML as 
		  (select xmlparse(document textdata) as xmlData from bulkLoadText)
		select 
		  xmlData,
		  (xpath('/ns:ListBucketResult/ns:Name/text()', xmlData, awsxmlns()))[1]::text,
		  (xpath('/ns:ListBucketResult/ns:Prefix/text()', xmlData, awsxmlns()))[1]::text,
		  (xpath('/ns:ListBucketResult/ns:IsTruncated/text()', xmlData, awsxmlns()))[1]::text::bool
			from parsedXML;


	insert into AWSBulkLoadTable
	           (bucket, prefix, isTruncated)
		select bucket, prefix, isTruncated
		from bulkLoadXML
		where not exists 
		(select 1 from AWSBulkLoadTable
		  where bucket = bulkLoadXML.bucket
		    and prefix = bulkLoadXML.prefix);

	select * from AWSBulkLoadTable;
		      
	update AWSBulkLoadTable
		set bucket = bulkLoadXML.bucket,
		    prefix = bulkLoadXML.prefix,
		    isTruncated = bulkLoadXML.isTruncated
		from bulkLoadXML
		where AWSBulkLoadTable.bucket = bulkLoadXML.bucket
		  and AWSBulkLoadTable.prefix = bulkLoadXML.prefix;


	insert into AWSBulkLoadDataTable
	  (loadID, path, timestamp, bytes, storage)
	select
	  AWSBulkLoadTable.loadID,
	  ((xpath('Key/text()', contents))[1]::text),
	  ((xpath('LastModified/text()', contents))[1]::text::timestamptz),
	  ((xpath('Size/text()', contents))[1]::text::integer),
	  ((xpath('StorageClass/text()', contents))[1]::text)
	from (select unnest(xpath('/ns:ListBucketResult/ns:Contents', xmlData, namespace())) as contents
	    from (select xmlData, bucket, prefix from bulkLoadXML) x1) x2
	  join AWSBulkLoadTable on
	    AWSBulkLoadTable.bucket = bucket and 
		AWSBulkLoadTable.prefix = prefix;


	end;
	$$
	language plpgsql
	returns null on null input;


-- Statistics:


create function humanReadableBytes(size bigint) returns text as
	$$
	declare 
	num double precision;
	unit text;
	begin
	
	if size is null then 
		return null;
		end if;
	
	num = size::double precision;
	foreach unit in array array[' bytes',' KB',' MB',' GB',' TB',' PB'] loop
		if num < 1024.0 then
			return to_char(num, '9999D9')||unit;
			end if;
		num = num / 1024.0;
		end loop;

	return num||' Really?';
	
	end; 
	$$
	language plpgsql immutable
	returns null on null input;


create function daySpan(fromDate timestamptz, toDate timestamptz) returns integer as
	$$
	declare quanta interval;
	declare result integer;
	begin

	if fromDate is null or toDate is null then
		return null;
		end if;
		
	quanta = fromDate - toDate;
	result = extract(epoch from quanta)::integer;
	return (result / (60 * 60 * 24))::integer;
	
	end;
	$$ 
	language plpgsql immutable
	returns null on null input;


create materialized view AWSObjectTableTotals as
select
	split_part(key, '/', 1) as Bundle,
	sum(case when awsdate < localdate then 1 else 0 end) as UpdatedParts,
	sum(case when awsdate < localdate then localbytes else 0 end) as UpdatedBytes,
	sum(case when awsdate is null then 1 else 0 end) as NewParts,
	sum(case when awsdate is null then localbytes else 0 end) as NewBytes,
	sum(case when localdate is null then 1 else 0 end) as DeleteParts,
	sum(case when localdate is null then awsbytes else 0 end) as DeleteBytes,
	sum(case when localdate is not null then 1 else 0 end) as TotalLocalParts,
	sum(case when localdate is not null then localbytes else 0 end) as TotalLocalBytes,
	sum(case when awsdate is not null then 1 else 0 end) as TotalAWSParts,
	sum(case when awsdate is not null then awsbytes else 0 end) as TotalAWSBytes,
	sum(case when awsstorage = 'AWSStorageGlacier' then 1 else 0 end) as totalGlacierParts,
	sum(case when awsstorage = 'AWSStorageGlacier' then awsbytes else 0 end) as totalGlacierBytes
		from awsobjecttable
		group by split_part(key, '/', 1);


create materialized view AWSObjectDistributionTable as
select
    awsDate::date as "Day",	
	split_part(key, '/', 1) as Bundle,
	sum(awsBytes::bigint) as "Bytes"
    from AWSObjectTable
	where awsDate is not null
		group by awsDate::date, split_part(key, '/', 1)
		order by awsDate::date desc;

-- date_trunc('day', awsDate)	

--select
--	split_part(key, '/', 1) as Bundle,
--		from aws

