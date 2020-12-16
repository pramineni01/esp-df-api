"""Create DF Tables

Revision ID: 15a34f49ff46
Revises: 
Create Date: 2020-06-30 12:00:18.962343

"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = '15a34f49ff46'
down_revision = None
branch_labels = None
depends_on = None


def upgrade():
    forecasts_sql="""
            create table forecasts(
            forecast_id INT UNSIGNED NOT NULL AUTO_INCREMENT,		 
            dataset_id BIGINT UNSIGNED NOT NULL,
            latest_version_dimension_member_id INT UNSIGNED NULL,
            start_timestamp TIMESTAMP(6) GENERATED ALWAYS AS ROW START INVISIBLE,		 
            end_timestamp TIMESTAMP(6) GENERATED ALWAYS AS ROW END INVISIBLE,
            PERIOD FOR SYSTEM_TIME(start_timestamp, end_timestamp),
            PRIMARY KEY (forecast_id,end_timestamp),
            UNIQUE INDEX forecasts_idx USING BTREE  (dataset_id,end_timestamp)  )ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 WITH SYSTEM VERSIONING
    """
    op.execute(forecasts_sql)

    forecast_translations_sql="""	
            create table forecast_translations(
            locale_id SMALLINT UNSIGNED NOT NULL,
            forecast_id INT UNSIGNED NOT NULL,		 
            forecast_name VARCHAR(256) NOT NULL,
            start_timestamp TIMESTAMP(6) GENERATED ALWAYS AS ROW START INVISIBLE,		 
            end_timestamp TIMESTAMP(6) GENERATED ALWAYS AS ROW END INVISIBLE,
            PERIOD FOR SYSTEM_TIME(start_timestamp, end_timestamp),
            PRIMARY KEY (locale_id,forecast_id,end_timestamp) )ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 WITH SYSTEM VERSIONING
    """
    op.execute(forecast_translations_sql)
    
    scenarios_sql="""
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
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 WITH SYSTEM VERSIONING


    """
    op.execute(scenarios_sql)
    
    scenario_tags_sql="""
            CREATE TABLE scenario_tags (
              scenario_id int(10) unsigned NOT NULL AUTO_INCREMENT,
              tag_id int(10) unsigned NOT NULL,
              start_timestamp timestamp(6) GENERATED ALWAYS AS ROW START INVISIBLE,
              end_timestamp timestamp(6) GENERATED ALWAYS AS ROW END INVISIBLE,
              PRIMARY KEY (scenario_id,tag_id,end_timestamp),
              KEY scenario_tags_scenario_id_idx (scenario_id,end_timestamp) USING BTREE,
              KEY scenario_tags_tag_id_idx (tag_id,end_timestamp) USING BTREE,
              PERIOD FOR SYSTEM_TIME (start_timestamp, end_timestamp)
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 WITH SYSTEM VERSIONING

    """
    
    op.execute(scenario_tags_sql)
    
    scenario_comments_sql="""	
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
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 WITH SYSTEM VERSIONING


    """
    op.execute(scenario_comments_sql)
    
    
    scenario_runs_sql="""		
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
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 WITH SYSTEM VERSIONING

    """
    
    op.execute(scenario_runs_sql)
    
    datasets_sql="""
            create table datasets(
            dataset_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
            datasource_id INT UNSIGNED NOT NULL,
            dataset_name VARCHAR(50) NOT NULL,
            dataset_description VARCHAR(1024) NULL,
            dataset_version VARCHAR(50) NOT NULL,
            PERIOD FOR SYSTEM_TIME(start_timestamp, end_timestamp),
            start_timestamp TIMESTAMP(6) GENERATED ALWAYS AS ROW START INVISIBLE,		 
            end_timestamp TIMESTAMP(6) GENERATED ALWAYS AS ROW END INVISIBLE,
            PRIMARY KEY (dataset_id,end_timestamp))ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 WITH SYSTEM VERSIONING
    """
    op.execute(datasets_sql)

    data_filters_sql="""
            create table data_filters (
            data_filter_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
            data_filter_name VARCHAR(256) NOT NULL,
            data_filter_definition JSON,
            user_id VARCHAR(36)  NULL,
            PERIOD FOR SYSTEM_TIME(start_timestamp, end_timestamp),
            start_timestamp TIMESTAMP(6) GENERATED ALWAYS AS ROW START INVISIBLE,
            end_timestamp TIMESTAMP(6) GENERATED ALWAYS AS ROW END INVISIBLE,
            PRIMARY KEY (data_filter_id, end_timestamp))ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 WITH SYSTEM VERSIONING
    """
    op.execute(data_filters_sql)

def downgrade():
    op.execute("""DROP TABLE IF EXISTS forecasts """ )
    op.execute("""DROP TABLE IF EXISTS forecast_translations """)    
    op.execute("""DROP TABLE IF EXISTS scenarios """)    
    op.execute("""DROP TABLE IF EXISTS scenario_tags """)    
    op.execute("""DROP TABLE IF EXISTS scenario_comments """)    
    op.execute("""DROP TABLE IF EXISTS scenario_runs """)
    op.execute("""DROP TABLE IF EXISTS datasets """)
    op.execute("""DROP TABLE IF EXISTS data_filters """)
