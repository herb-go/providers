create table if not exists cache(
	        cache_name varchar(255),
            cache_key varchar(767),
            version varchar(255),
	        cache_value LONGBLOB,
	        expired bigint,
	        primary key (cache_name,cache_key)
        ) ENGINE=InnoDB charset latin1 COLLATE latin1_bin;