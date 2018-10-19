create table if not exists cache(
	        cache_name varchar(255),
            cache_key varchar(255),
            version varchar(255),
	        cache_value BLOB,
	        expired bigint,
	        primary key (cache_name,cache_key)
        ) 