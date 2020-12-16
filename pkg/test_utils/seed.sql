-- Data models test

INSERT INTO datasets
            (dataset_id, datasource_id, dataset_name, dataset_description, dataset_version)
VALUES
  (1, 1, "test_dataset", "data set for test", 0.1);


INSERT INTO scenarios
            (scenario_id, forecast_id, scenario_name, branch_id, user_id, impute_history, scenario_status)
VALUES
  (1, 2, "delete_test1", 2, 0, true, "CURRENT");

INSERT INTO scenarios
            (scenario_id, forecast_id, scenario_name, branch_id, user_id, impute_history, scenario_status)
VALUES
  (2, 2, "delete_test2", 2, 0, true, "CURRENT");


INSERT INTO scenarios
            (scenario_id, forecast_id, scenario_name, branch_id, user_id, impute_history, scenario_status)
VALUES
  (2, 2, "delete_test2", 2, 0, true, "CURRENT");


INSERT INTO scenario_comments
            (scenario_comment_id, scenario_id, comment)
VALUES
  (1, 1, "comment1");

INSERT INTO scenario_comments
            (scenario_comment_id, scenario_id, comment)
VALUES
  (2, 1, "comment2");

INSERT INTO scenario_comments
            (scenario_comment_id, scenario_id, comment)
VALUES
  (3, 2, "comment3");
