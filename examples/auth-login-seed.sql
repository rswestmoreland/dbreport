INSERT INTO user_accounts (id, full_name, username, email) VALUES
(1,'Avery Holt','aholt','avery.holt@example.invalid'),
(2,'Jordan Pike','jpike','jordan.pike@test.test'),
(3,'Taylor Wynn','twynn','taylor.wynn@corp.example'),
(4,'Morgan Vale','mvale','morgan.vale@example.invalid'),
(5,'Riley Shaw','rshaw','riley.shaw@test.test'),
(6,'Casey Dean','cdean','casey.dean@corp.example'),
(7,'Quinn Hale','qhale','quinn.hale@example.invalid'),
(8,'Parker Drew','pdrew','parker.drew@test.test'),
(9,'Skyler Boone','sboone','skyler.boone@corp.example'),
(10,'Emery Knox','eknox','emery.knox@example.invalid'),
(11,'Rowan Blair','rblair','rowan.blair@test.test'),
(12,'Harper Lane','hlane','harper.lane@corp.example'),
(13,'Finley Reeve','freeve','finley.reeve@example.invalid'),
(14,'Kendall Moss','kmoss','kendall.moss@test.test'),
(15,'Sage Pruitt','spruitt','sage.pruitt@corp.example'),
(16,'Dakota Flynn','dflynn','dakota.flynn@example.invalid'),
(17,'Cameron Bell','cbell','cameron.bell@test.test'),
(18,'Reese Monroe','rmonroe','reese.monroe@corp.example'),
(19,'Alex Mercer','amercer','alex.mercer@example.invalid'),
(20,'Bailey Stone','bstone','bailey.stone@test.test'),
(21,'Phoenix Clarke','pclarke','phoenix.clarke@corp.example'),
(22,'Jules Hart','jhart','jules.hart@example.invalid'),
(23,'Remy Cross','rcross','remy.cross@test.test'),
(24,'Noel Quinn','nquinn','noel.quinn@corp.example');

INSERT INTO auth_attempts (id, user_id, login_time, result, reason)
WITH RECURSIVE seq(n) AS (
  SELECT 1 UNION ALL SELECT n + 1 FROM seq WHERE n < 210
)
SELECT
  n,
  ((n - 1) % 24) + 1,
  TIMESTAMP('2026-02-01 00:00:00') + INTERVAL (n - 1) HOUR,
  CASE WHEN MOD(n, 7) IN (0, 3) THEN 'failure' ELSE 'success' END,
  CASE
    WHEN MOD(n, 7) IN (0, 3) THEN ELT((MOD(n, 3) + 1), 'bad password', 'non-existent username', 'locked account')
    ELSE 'authenticated'
  END
FROM seq;

INSERT INTO user_agent_info (id, auth_attempt_id, browser, country)
SELECT
  id,
  id,
  CASE
    WHEN MOD(id, 2) = 0 THEN 'Chrome'
    WHEN MOD(id, 5) = 0 THEN 'Safari'
    WHEN MOD(id, 7) = 0 THEN 'Firefox'
    WHEN MOD(id, 11) = 0 THEN 'Edge'
    WHEN MOD(id, 13) = 0 THEN 'Brave'
    ELSE 'Opera'
  END,
  CASE
    WHEN MOD(id, 2) = 0 THEN 'Australia'
    WHEN MOD(id, 3) = 0 THEN 'Canada'
    WHEN MOD(id, 5) = 0 THEN 'Japan'
    WHEN MOD(id, 7) = 0 THEN 'Germany'
    WHEN MOD(id, 11) = 0 THEN 'India'
    WHEN MOD(id, 13) = 0 THEN 'Brazil'
    ELSE 'United States'
  END
FROM auth_attempts;
