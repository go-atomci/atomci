# support normal docker registry
UPDATE `sys_integrate_setting` SET `type`='registry' WHERE `type`='harbor';

ALTER TABLE `project_env` CHANGE COLUMN `harbor` `registry` bigint(20) NOT NULL DEFAULT 0 AFTER `ci_server`;