begin;
truncate table bulkLoadText cascade;
truncate table bulkLoadXML cascade;
truncate table AWSBulkLoadTable cascade;
truncate table AWSBulkLoadDataTable cascade;
commit;

begin;

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


drop function if exists namespace() cascade;
create function namespace() returns text[][] as
	$$
	begin
	return array[array['ns', 'http://s3.amazonaws.com/doc/2006-03-01/']];
	end;
	$$ 
	language plpgsql immutable;


select namespace();


insert into bulkLoadXML
	with parsedXML as 
	  (select xmlparse(document textdata) as xmlData from bulkLoadText)
	select 
	  xmlData,
	  (xpath('/ns:ListBucketResult/ns:Name/text()', xmlData, namespace()))[1]::text,
	  (xpath('/ns:ListBucketResult/ns:Prefix/text()', xmlData, namespace()))[1]::text,
	  (xpath('/ns:ListBucketResult/ns:IsTruncated/text()', xmlData, namespace()))[1]::text::bool
		from parsedXML;

-- select * from bulkLoadXML;

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


commit;

select * from AWSBulkLoadDataTable;

-- select path from awsbulkloaddataTable where loadID = 28 order by path desc  limit 1;
select path from awsbulkloaddataTable order by path desc  limit 1;
