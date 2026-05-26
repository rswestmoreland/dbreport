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
WITH RECURSIVE seq(n) AS (SELECT 1 UNION ALL SELECT n + 1 FROM seq WHERE n < 233),
plan(day, success_total, failure_total) AS (
  SELECT DATE('2026-02-01'), 28, 4 UNION ALL
  SELECT DATE('2026-02-02'), 18, 3 UNION ALL
  SELECT DATE('2026-02-03'), 35, 6 UNION ALL
  SELECT DATE('2026-02-04'), 22, 5 UNION ALL
  SELECT DATE('2026-02-05'), 24, 31 UNION ALL
  SELECT DATE('2026-02-06'), 30, 7 UNION ALL
  SELECT DATE('2026-02-07'), 16, 4
), daily AS (
  SELECT day, success_total, failure_total, success_total + failure_total AS total,
         SUM(success_total + failure_total) OVER (ORDER BY day) AS cumulative_total
  FROM plan
)
SELECT
  seq.n,
  ((seq.n - 1) % 24) + 1,
  TIMESTAMP(d.day) + INTERVAL MOD(seq.n*37, 24) HOUR + INTERVAL MOD(seq.n*13, 60) MINUTE + INTERVAL MOD(seq.n*17, 60) SECOND,
  CASE WHEN (seq.n - d.start_n + 1) <= d.failure_total THEN 'failure' ELSE 'success' END,
  CASE
    WHEN (seq.n - d.start_n + 1) <= d.failure_total THEN ELT((MOD(seq.n, 3) + 1), 'bad password', 'non-existent username', 'locked account')
    ELSE 'authenticated'
  END
FROM seq
JOIN (
  SELECT day, success_total, failure_total, total, cumulative_total,
         (cumulative_total - total + 1) AS start_n
  FROM daily
) d ON seq.n BETWEEN d.start_n AND d.cumulative_total;

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
