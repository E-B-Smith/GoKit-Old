
begin;
create temp table AWSBulkLoadText (textData text) on commit drop;
copy bucketListText (textData)
	from '/Users/Edward/Development/go/src/violent.blue/go/aws-bu/TestData/TestBackup.adata'
	from '/Users/Edward/Development/aws-backup/AWSUploadTest/TestBackup.xml';
--	from stdin;

-- Create a temp table to hold the parsed XML so we only parse it once:

create temp table bulkLoadXML
	(
	xmlData  	xml,
	bucket	 	varchar(256),
	prefix	 	varchar(512),
	isTruncated bool
	) 
	on commit drop;

insert into AWSBulkLoadXML
	with parsedXML as 
	  (select xmlparse(document textdata) as xmlData from AWSBulkLoadText)
	select 
	  xmlData,
	  (xpath('/ListBucketResult/Name/text()', xmlData))[1]::text,
	  (xpath('/ListBucketResult/Prefix/text()', xmlData))[1]::text,
	  (xpath('/ListBucketResult/Marker/text()', xmlData))[1]::text,
	  (xpath('/ListBucketResult/IsTruncated/text()', xmlData))[1]::text::bool
		from parsedXML;

insert into AWSBulkLoadTable
           (bucket, prefix, isTruncated)
	select bucket, prefix, isTruncated
	from bulkLoadXML
	where not exists 
	(select 1 from AWSBulkLoadTable
	  where bucket = AWSBulkLoadXML.bucket and 
	        prefix = AWSBulkLoadXML.prefix);
	      
update AWSBulkLoadTable
	set bucket = AWSBulkLoadXML.bucket,
	    prefix = AWSBulkLoadXML.prefix,
	    isTruncated = AWSBulkLoadXML.isTruncated
	from bucketListXML
	where AWSBulkLoadTable.bucket = bucketListXML.bucket
	  and AWSBulkLoadTable.prefix = bucketListXML.prefix;

insert into AWSBulkLoadDataTable
  (bucketListID, path, mdate, state)
select
  AWSBulkLoadTable.loadID,
  ((xpath('Key/text()', contents))[1]::text),
  ((xpath('LastModified/text()', contents))[1]::text::timestamptz),
  ((xpath('StorageClass/text()', contents))[1]::text)
from (select unnest(xpath('/ListBucketResult/Contents', xmlData)) as contents
    from (select xmlData, bucket, prefix from bucketListXML) x1) x2
  join awsBucketListTable on
    awsBucketListTable.bucket = bucket and 
	awsBucketListTable.prefix = prefix;

commit;
