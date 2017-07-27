-- create table
CREATE TABLE IF NOT EXISTS `nested`(
`id` BIGINT NOT NULL COMMENT 'node ID',
`node` VARCHAR(64) CHARACTER SET 'utf8' NOT NULL COMMENT 'node name',
`pid` BIGINT NOT NULL COMMENT 'parent ID',
`depth` INT NOT NULL COMMENT 'Level',
`lft` INT NOT NULL COMMENT 'left index',
`rgt` INT NOT NULL COMMENT 'right index',
  PRIMARY KEY (`id`),
  INDEX `depth_index` (`depth` ASC),
  INDEX `lft_index` (`lft` ASC),
  INDEX `rgt_index` (`rgt` ASC))
ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COMMENT = 'nested sets model';
