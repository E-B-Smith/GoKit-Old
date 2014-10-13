
drop database if exists AWSObjects;
create database AWSObjects with encoding 'UTF8';


drop table if exists AWSParameterTable;
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


drop type if exists AWSStorageState cascade;
create type AWSStorageState as enum
	(
	 'AWSStorageLocal'
	,'AWSStorageGlacier'
	,'AWSStorageRestoring'
	,'AWSStorageRestored'
	,'AWSStorageStandard'
	);

drop table if exists AWSObjectTable cascade;
create table AWSObjectTable
	(
	 key				varchar(256) unique not null primary key

	,awsDate			timestamptz
	,awsBytes			integer
	,awsStorage			AWSStorageState

	,localDate			timestamptz
	,localBytes			integer
	);

drop type if exists AWSLogLevel cascade;
create type AWSLogLevel as enum
	(
	 'AWSLogDebug'
	,'AWSLogInfo'
	,'AWSLogStart'
	,'AWSLogExit'
	,'AWSLogWarning'
	,'AWSLogError'
	);

drop table if exists AWSLogTable;
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

drop index if exists AWSLogTimeIndex;
create index AWSLogTimeIndex on AWSLogTable(time, entry);


drop table if exists AWSStatusTable;
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

drop index if exists AWSStatusPIDIndex;
create index AWSStatusPIDIndex on AWSStatusTable(pid);



-- Bulk load tables:



drop table if exists AWSBulkLoadTable cascade;
create table AWSBulkLoadTable
	(
	loadID		   serial primary key,
	bucket         varchar(255) not null,
	prefix		   varchar(512) not null,
	isTruncated    bool         not null,
	
	constraint AWSBulkLoadTableUniqueConstraint
		unique (bucket, prefix)
	);


drop table if exists AWSBulkLoadDataTable; 
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


drop index if exists AWSBulkLoadDataPathIndex;
create index AWSBuildLoadDataPathIndex on AWSBulkLoadDataTable(loadID, path);



-- Statistics:



drop materialized view if exists AWSObjectTableTotals;
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


-- drop materialized view if exists AWSObjectDistributionTable;
-- create materialized view AWSObjectDistributionTable as
-- select
--     awsDate::date as "Day",	
-- 	split_part(key, '/', 1) as Bundle,
-- 	sum(awsBytes::bigint) as "Bytes"
--     from AWSObjectTable
-- 	where awsDate is not null
-- 		group by awsDate::date 
-- 		order by awsDate::date desc;

-- date_trunc('day', awsDate)	

--select
--	split_part(key, '/', 1) as Bundle,
--		from aws


drop function if exists humanReadableBytes(size bigint) cascade;
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



drop function if exists daySpan(fromDate timestamptz, toDate timestamptz) cascade;
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


