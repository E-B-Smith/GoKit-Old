begin;

lock table AWSObjectTable in exclusive mode;

update AWSObjectTable
	set awsDate = AWSBulkLoadDataTable.awsDate,
		awsBytes = AWSBulkLoadDataTable.awsBytes,
		awsStorage = AWSBulkLoadDataTable.awsStorage
	from AWSBulkLoadDataTable
	where AWSBulkLoadDataTable.key = AWSObjectTable.key;
	
insert into AWSObjectTable
		(key, awsDate, awsBytes, awsStorage)
	select
		AWSBulkLoadDataTable.key,
		AWSBulkLoadDataTable.awsDate,
		AWSBulkLoadDataTable.awsBytes,
		AWSBulkLoadDataTable.awsStorage,
	from AWSBulkLoadDataTable
	left outer join AWSObjectTable on 
		(AWSObjectTable.key = AWSBulkLoadDataTable.key)
	where AWSObjectTable.key is null;
	
commit;

select * from AWSObjectTable;
