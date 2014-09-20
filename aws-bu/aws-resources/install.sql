
drop database if exists AWSObjects;
create database AWSObjects with encoding 'UTF8';

drop type if exists AWSObjectState cascade;
create type AWSObjectState as enum
	(
	 'AWSStateGlacier'
	,'AWSStateRestoring'
	,'AWSStateRestored'
	,'AWSStateStandard'
	);

drop table if exists AWSObjectTable;
create table AWSObjectTable
	(
	 key				varchar(256) unique not null primary key

	,awsDate			timestamptz
	,awsBytes			integer
	,awsState			AWSObjectState

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
	sum(case when awsstate = 'AWSStateGlacier' then 1 else 0 end) as totalGlacierParts,
	sum(case when awsstate = 'AWSStateGlacier' then awsbytes else 0 end) as totalGlacierBytes
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

	return lpad('Really?', 8);
	
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


drop table if exists AWSParameterTable;
create table AWSParameterTable
	(
	 version 			varchar(10)
	,bundlepath			varchar(256)
	,AWSAccessKeyID		varchar(64)
	,AWSAccessSecret	varchar(64)
	,AWSRegion			varchar(16)
	);

insert into AWSParameterTable (version) values ('1.00.001');

