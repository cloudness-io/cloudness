CREATE TABLE servers (
 server_id									INTEGER PRIMARY KEY AUTOINCREMENT
,server_uid									INTEGER NOT NULL
,server_type								TEXT NOT NULL
,server_name								TEXT NOT NULL
,server_description							TEXT
,server_ipv4								TEXT
,server_ipv6								TEXT
,server_wildcard_domain						TEXT
,server_dns_proxy							TEXT
,server_proxy_auth_key						TEXT
,server_user								TEXT
,server_port								INTEGER 
,server_volume_supports_online_expansion     BOOLEAN
,server_builder_is_enabled					BOOLEAN
,server_builder_is_build_server			    BOOLEAN
,server_builder_max_concurrent_builds	    INTEGER
,server_builder_max_cpu						REAL	
,server_builder_max_memory					REAL	
,server_volume_min_size                     INTEGER
,server_created								INTEGER
,server_updated								INTEGER
);
