create database AWSObjects with encoding 'UTF8';
-- create user AWSBackup;

create type AWSObjectState as enum
	(
	 'AWSStateGlacier'
	,'AWSStateRestoring'
	,'AWSStateRestored'
	,'AWSStateStandard'
	);

drop table AWSObjectTable;
create table AWSObjectTable
	(
	 key				varchar(256) unique not null primary key

	,awsDate			timestamptz
	,awsBytes			integer
	,awsState			AWSObjectState

	,localDate			timestamptz
	,localBytes			integer
	);

drop index AWSObjectKeyIndex;
create unique index AWSObjectKeyIndex on AWSObjectTable(key);
	
create type AWSAction as enum
	(
	 'AWSActionBackup'
	,'AWSActionRestore'
	,'AWSActionRefresh'
	);

drop table AWSActionTable;
create table AWSActionTable
	(
	 startDate		timestamptz		unique not null primary key
	,finishDate		timestamptz
	,pid 			integer
	,exitStatus		integer
	,action			AWSAction		not null
	,keyPrefix		varchar(255)	not null
	,objects		integer
	,bytes			integer
	,totalObjects	integer
	,totalBytes		integer
	);

drop index AWSActionDateIndex;
create unique index AWSActionDateIndex on AWSActionTable(startDate);

create type AWSLogLevel as enum
	(
	 'AWSLogDebug'
	,'AWSLogInfo'
	,'AWSLogWarning'
	,'AWSLogError'
	);

drop table AWSLogTable;
create table AWSLogTable
	(
	 time		timestamptz		unique not null primary key
	,pid 		integer
	,level		AWSLogLevel
	,message	varchar(512)
	);

drop index AWSLogIndex;
create unique index AWSLogIndex on AWSLogTable(time);


create or replace function humanReadableBytes(size bigint) returns text as
	\$\$
	declare 
	num double precision;
	unit text;
	begin
	
	num = size::double precision;
	foreach unit in array array[' bytes',' KB',' MB',' GB',' TB',' PB',' Really?'] loop
		if num < 1024.0 then
			return to_char(num, '9999D9')||unit;
			end if;
		num = num / 1024.0;
		end loop;
		
	end; 
	\$\$ 
	language plpgsql immutable;


create or replace function hoursSince(startTime timestamptz, timeIn timestamptz) returns integer as
	\$\$
	declare quanta interval;
	declare result integer;
	begin

	quanta = timeIn - startTime;
	result = extract(epoch from quanta)::integer;
	
	return (result / (60 * 60 * 24))::integer;
	
	end;
	\$\$ 
	language plpgsql immutable;

