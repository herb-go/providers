create table if not exists cache(
	        cache_name varchar(255),
            cache_key varchar(255),
            version varchar(255),
	        cache_value Text,
	        expired bigint,
	        primary key (cache_name,cache_key)
        ) 