CREATE TABLE alembic_version (
    version_num VARCHAR(32) NOT NULL, 
    CONSTRAINT alembic_version_pkc PRIMARY KEY (version_num)
);

-- Running upgrade  -> 15a34f49ff46

create table forecasts(
            forecast_id INT UNSIGNED NOT NULL AUTO_INCREMENT,         
            dataset_id BIGINT UNSIGNED NOT NULL,
            latest_version_dimension_member_id INT UNSIGNED NULL,
            start_timestamp TIMESTAMP(6) GENERATED ALWAYS AS ROW START INVISIBLE,         
            end_timestamp TIMESTAMP(6) GENERATED ALWAYS AS ROW END INVISIBLE,
            PERIOD FOR SYSTEM_TIME(start_timestamp, end_timestamp),
            PRIMARY KEY (forecast_id,end_timestamp),
            UNIQUE INDEX forecasts_idx USING BTREE  (dataset_id,end_timestamp)  )ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 WITH SYSTEM VERSIONING;

create table forecast_translations(
            locale_id SMALLINT UNSIGNED NOT NULL,
            forecast_id INT UNSIGNED NOT NULL,         
            forecast_name VARCHAR(256) NOT NULL,
            start_timestamp TIMESTAMP(6) GENERATED ALWAYS AS ROW START INVISIBLE,         
            end_timestamp TIMESTAMP(6) GENERATED ALWAYS AS ROW END INVISIBLE,
            PERIOD FOR SYSTEM_TIME(start_timestamp, end_timestamp),
            PRIMARY KEY (locale_id,forecast_id,end_timestamp) )ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 WITH SYSTEM VERSIONING;

CREATE TABLE scenarios (
              scenario_id int(10) unsigned NOT NULL AUTO_INCREMENT,
              forecast_id int(10) unsigned NOT NULL,
              user_id varchar(36) DEFAULT NULL,
              scenario_name varchar(1024) NOT NULL,
              scope_id varchar(512) NOT NULL,
              da_branch_id int(10) unsigned NOT NULL,
              scenario_status enum('CURRENT','DELETED','PROMOTED','SUPERSCEDED') NOT NULL DEFAULT 'CURRENT',
              start_timestamp timestamp(6) GENERATED ALWAYS AS ROW START INVISIBLE,
              end_timestamp timestamp(6) GENERATED ALWAYS AS ROW END INVISIBLE,
              PRIMARY KEY (scenario_id,end_timestamp),
              KEY scenarios_user_idx (user_id,end_timestamp) USING BTREE,
              KEY scenarios_forecast_idx (forecast_id,end_timestamp) USING BTREE,
              KEY scenarios_scope_idx (scope_id,end_timestamp) USING BTREE,
              KEY scenarios_scenario_status_idx (scenario_status,end_timestamp) USING BTREE,
              PERIOD FOR SYSTEM_TIME (start_timestamp, end_timestamp)
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 WITH SYSTEM VERSIONING;

CREATE TABLE scenario_tags (
              scenario_id int(10) unsigned NOT NULL AUTO_INCREMENT,
              tag_id int(10) unsigned NOT NULL,
              start_timestamp timestamp(6) GENERATED ALWAYS AS ROW START INVISIBLE,
              end_timestamp timestamp(6) GENERATED ALWAYS AS ROW END INVISIBLE,
              PRIMARY KEY (scenario_id,tag_id,end_timestamp),
              KEY scenario_tags_scenario_id_idx (scenario_id,end_timestamp) USING BTREE,
              KEY scenario_tags_tag_id_idx (tag_id,end_timestamp) USING BTREE,
              PERIOD FOR SYSTEM_TIME (start_timestamp, end_timestamp)
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 WITH SYSTEM VERSIONING;

CREATE TABLE scenario_comments (
          scenario_comment_id int(10) unsigned NOT NULL AUTO_INCREMENT,
          scenario_id int(10) unsigned NOT NULL,
          comment varchar(1024) NOT NULL,
          start_timestamp timestamp(6) GENERATED ALWAYS AS ROW START INVISIBLE,
          end_timestamp timestamp(6) GENERATED ALWAYS AS ROW END INVISIBLE,
          user_id varchar(36) NOT NULL,
          PRIMARY KEY (scenario_comment_id,end_timestamp),
          KEY scenario_comments_scenario_id_idx (scenario_id,end_timestamp) USING BTREE,
          PERIOD FOR SYSTEM_TIME (start_timestamp, end_timestamp)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 WITH SYSTEM VERSIONING;

CREATE TABLE scenario_runs (
              scenario_run_id int(10) unsigned NOT NULL AUTO_INCREMENT,
              scenario_id int(10) unsigned NOT NULL,
              user_id varchar(36) DEFAULT NULL,
              scenario_run_status enum('SCHEDULED','IN_PROGRESS','FORECASTED','ERROR') NOT NULL DEFAULT 'SCHEDULED',
              run_start_timestamp timestamp(6) NULL DEFAULT NULL,
              run_end_timestamp timestamp(6) NULL DEFAULT NULL,
              da_version_id bigint(20) unsigned DEFAULT NULL,
              start_timestamp timestamp(6) GENERATED ALWAYS AS ROW START INVISIBLE,
              end_timestamp timestamp(6) GENERATED ALWAYS AS ROW END INVISIBLE,
              PRIMARY KEY (scenario_run_id,end_timestamp),
              KEY scenarios_scenario_idx (scenario_id,end_timestamp) USING BTREE,
              KEY scenarios_user_idx (user_id,end_timestamp) USING BTREE,
              PERIOD FOR SYSTEM_TIME (start_timestamp, end_timestamp)
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 WITH SYSTEM VERSIONING;

create table datasets(
            dataset_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
            datasource_id INT UNSIGNED NOT NULL,
            dataset_name VARCHAR(50) NOT NULL,
            dataset_description VARCHAR(1024) NULL,
            dataset_version VARCHAR(50) NOT NULL,
            PERIOD FOR SYSTEM_TIME(start_timestamp, end_timestamp),
            start_timestamp TIMESTAMP(6) GENERATED ALWAYS AS ROW START INVISIBLE,         
            end_timestamp TIMESTAMP(6) GENERATED ALWAYS AS ROW END INVISIBLE,
            PRIMARY KEY (dataset_id,end_timestamp))ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 WITH SYSTEM VERSIONING;

create table data_filters (
            data_filter_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
            data_filter_name VARCHAR(256) NOT NULL,
            data_filter_definition JSON,
            user_id VARCHAR(36)  NULL,
            PERIOD FOR SYSTEM_TIME(start_timestamp, end_timestamp),
            start_timestamp TIMESTAMP(6) GENERATED ALWAYS AS ROW START INVISIBLE,
            end_timestamp TIMESTAMP(6) GENERATED ALWAYS AS ROW END INVISIBLE,
            PRIMARY KEY (data_filter_id, end_timestamp))ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 WITH SYSTEM VERSIONING;

INSERT INTO alembic_version (version_num) VALUES ('15a34f49ff46');

