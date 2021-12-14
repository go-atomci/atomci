ALTER TABLE atomci.sys_group_role_operation DROP KEY `group`;

ALTER TABLE atomci.sys_group_role_operation ADD CONSTRAINT roleOperation UNIQUE (role, operation_id);