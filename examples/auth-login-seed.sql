INSERT INTO user_accounts (id, full_name, username, email) VALUES
(1,'Avery Holt','aholt','avery.holt@example.invalid'),
(2,'Jordan Pike','jpike','jordan.pike@test.invalid'),
(3,'Taylor Wynn','twynn','taylor.wynn@corp.invalid'),
(4,'Morgan Vale','mvale','morgan.vale@example.invalid'),
(5,'Riley Shaw','rshaw','riley.shaw@test.invalid'),
(6,'Casey Dean','cdean','casey.dean@corp.invalid'),
(7,'Quinn Hale','qhale','quinn.hale@example.invalid'),
(8,'Parker Drew','pdrew','parker.drew@test.invalid'),
(9,'Skyler Boone','sboone','skyler.boone@corp.invalid'),
(10,'Emery Knox','eknox','emery.knox@example.invalid'),
(11,'Rowan Blair','rblair','rowan.blair@test.invalid'),
(12,'Harper Lane','hlane','harper.lane@corp.invalid'),
(13,'Finley Reeve','freeve','finley.reeve@example.invalid'),
(14,'Kendall Moss','kmoss','kendall.moss@test.invalid'),
(15,'Sage Pruitt','spruitt','sage.pruitt@corp.invalid'),
(16,'Dakota Flynn','dflynn','dakota.flynn@example.invalid'),
(17,'Cameron Bell','cbell','cameron.bell@test.invalid'),
(18,'Reese Monroe','rmonroe','reese.monroe@corp.invalid'),
(19,'Alex Mercer','amercer','alex.mercer@example.invalid'),
(20,'Bailey Stone','bstone','bailey.stone@test.invalid'),
(21,'Phoenix Clarke','pclarke','phoenix.clarke@corp.invalid'),
(22,'Jules Hart','jhart','jules.hart@example.invalid'),
(23,'Remy Cross','rcross','remy.cross@test.invalid'),
(24,'Noel Quinn','nquinn','noel.quinn@corp.invalid');

INSERT INTO auth_attempts (id, user_id, login_time, result)
WITH RECURSIVE seq(n) AS (
  SELECT 1
  UNION ALL
  SELECT n + 1 FROM seq WHERE n < 240
)
SELECT
  n,
  ((n - 1) % 24) + 1,
  TIMESTAMP('2026-02-01 00:00:00') + INTERVAL (n * 3) HOUR,
  CASE WHEN MOD(n, 5) = 0 OR MOD(n, 11) = 0 THEN 'failure' ELSE 'success' END
FROM seq;

INSERT INTO user_agent_info (id, auth_attempt_id, browser, country)
SELECT
  id,
  id,
  ELT(MOD(id, 6) + 1, 'Chrome', 'Firefox', 'Safari', 'Edge', 'Brave', 'Opera'),
  ELT(MOD(id, 8) + 1, 'United States', 'Canada', 'Germany', 'Japan', 'Brazil', 'India', 'Australia', 'South Africa')
FROM auth_attempts;
