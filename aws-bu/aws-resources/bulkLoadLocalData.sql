begin;
truncate table AWSObjectTable;
truncate table BulkLoadTable cascade;
commit;

begin;

-- create temporary table BulkLoadTable
create table if not exists BulkLoadTable
	(
	 key				varchar(256) unique not null primary key
	,localDate			timestamptz
	,localBytes			integer
	)
; -- on commit drop;

copy BulkLoadTable 
	from '/Users/Edward/Development/go/src/violent.blue/go/aws-bu/TestData/TestBackup.ldata'
--	from stdin
	with null '';

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
	
commit;

select * from AWSObjectTable;
