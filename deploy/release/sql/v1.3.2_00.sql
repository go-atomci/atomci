# support normal docker registry
UPDATE `sys_integrate_setting` SET `type`='registry' WHERE `type`='harbor';


# modify table project_env column harbor to registry
DROP PROCEDURE IF EXISTS `ModifyHarborToRegistry`;
delimiter $$
CREATE PROCEDURE `ModifyHarborToRegistry`()
BEGIN
    DECLARE HARBOREXISTS int DEFAULT 0;
    DECLARE REGISTRYEXISTS int DEFAULT 0;
    SELECT count(1) INTO @HARBOREXISTS FROM information_schema.COLUMNS WHERE TABLE_NAME='project_env' AND COLUMN_NAME='harbor';
    SELECT count(1) INTO @REGISTRYEXISTS FROM information_schema.COLUMNS WHERE TABLE_NAME='project_env' AND COLUMN_NAME='registry';
    IF @HARBOREXISTS>0 AND @REGISTRYEXISTS=0 #存在harbor列 不存在registry列时 直接修改列名
    THEN
        ALTER TABLE `project_env` CHANGE COLUMN `harbor` `registry` bigint(20) NOT NULL DEFAULT 0;
    ELSEIF  @HARBOREXISTS>0 AND @REGISTRYEXISTS>0 #harbor列和registry都存在时迁移数据并删除harbor列
    THEN
        UPDATE `project_env` SET `registry`=`harbor`;
        ALTER TABLE `project_env` DROP COLUMN `harbor`;
    END IF;
END;
$$
delimiter ;
CALL `ModifyHarborToRegistry`;
DROP PROCEDURE IF EXISTS `ModifyHarborToRegistry`;