INSERT INTO user_accounts (id, full_name, username, email) VALUES
(1,'Avery Holt','aholt','avery.holt@meadow.invalid'),
(2,'Jordan Pike','jpike','jordan.pike@northstar.test'),
(3,'Taylor Wynn','twynn','taylor.wynn@acme.example'),
(4,'Morgan Vale','mvale','morgan.vale@meadow.invalid'),
(5,'Riley Shaw','rshaw','riley.shaw@northstar.test'),
(6,'Casey Dean','cdean','casey.dean@acme.example'),
(7,'Quinn Hale','qhale','quinn.hale@meadow.invalid'),
(8,'Parker Drew','pdrew','parker.drew@northstar.test'),
(9,'Skyler Boone','sboone','skyler.boone@acme.example'),
(10,'Emery Knox','eknox','emery.knox@meadow.invalid'),
(11,'Rowan Blair','rblair','rowan.blair@northstar.test'),
(12,'Harper Lane','hlane','harper.lane@acme.example'),
(13,'Finley Reeve','freeve','finley.reeve@meadow.invalid'),
(14,'Kendall Moss','kmoss','kendall.moss@northstar.test'),
(15,'Sage Pruitt','spruitt','sage.pruitt@acme.example'),
(16,'Dakota Flynn','dflynn','dakota.flynn@meadow.invalid'),
(17,'Cameron Bell','cbell','cameron.bell@northstar.test'),
(18,'Reese Monroe','rmonroe','reese.monroe@acme.example'),
(19,'Alex Mercer','amercer','alex.mercer@meadow.invalid'),
(20,'Bailey Stone','bstone','bailey.stone@northstar.test'),
(21,'Phoenix Clarke','pclarke','phoenix.clarke@acme.example'),
(22,'Jules Hart','jhart','jules.hart@meadow.invalid'),
(23,'Remy Cross','rcross','remy.cross@northstar.test'),
(24,'Noel Quinn','nquinn','noel.quinn@acme.example');

INSERT INTO auth_attempts (id, user_id, login_time, result, reason)
WITH RECURSIVE seq(n) AS (
  SELECT 1 UNION ALL SELECT n + 1 FROM seq WHERE n < 280
)
SELECT
  n,
  ((n - 1) % 24) + 1,
  TIMESTAMP('2026-02-01 08:07:13') + INTERVAL (n - 1) HOUR + INTERVAL MOD(n*7, 43) MINUTE + INTERVAL MOD(n*11, 53) SECOND,
  CASE
    WHEN DATE(TIMESTAMP('2026-02-01 08:07:13') + INTERVAL (n - 1) HOUR) = '2026-02-05' AND MOD(n, 3) != 0 THEN 'failure'
    WHEN MOD(n, 9) IN (0, 4, 7) THEN 'failure'
    ELSE 'success'
  END,
  CASE
    WHEN DATE(TIMESTAMP('2026-02-01 08:07:13') + INTERVAL (n - 1) HOUR) = '2026-02-05' AND MOD(n, 3) != 0 THEN ELT((MOD(n, 3) + 1), 'bad password', 'non-existent username', 'locked account')
    WHEN MOD(n, 9) IN (0, 4, 7) THEN ELT((MOD(n, 3) + 1), 'bad password', 'non-existent username', 'locked account')
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
