-- 创建数据库
CREATE DATABASE IF NOT EXISTS logistics_manage DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE logistics_manage;

-- 插入默认管理员账号 (密码: admin123)
INSERT INTO `admins` (`username`, `password`, `real_name`, `phone`, `role`, `status`, `created_at`, `updated_at`)
VALUES ('admin', '$2a$14$wVsaPvJnQPx3sXqFpQjG7uBn0uLwJ9kH5gF3dE2aS1cVbN9mK8', '超级管理员', '13800138000', 2, 1, NOW(), NOW());

-- 插入默认配送员账号 (密码: delivery123, 需要管理员审核)
INSERT INTO `delivery_men` (`username`, `password`, `phone`, `real_name`, `id_card`, `site_id`, `status`, `verify_status`, `created_at`, `updated_at`)
VALUES ('delivery', '$2a$14$wVsaPvJnQPx3sXqFpQjG7uBn0uLwJ9kH5gF3dE2aS1cVbN9mK8', '13800138001', '张配送', '110101199001011234', 1, 1, 1, NOW(), NOW());

-- 插入测试站点
INSERT INTO `sites` (`name`, `province`, `city`, `district`, `address`, `phone`, `manager`, `status`, `created_at`, `updated_at`)
VALUES 
('北京朝阳站点', '北京市', '北京市', '朝阳区', '北京市朝阳区建国路88号', '010-88888888', '王经理', 1, NOW(), NOW()),
('上海浦东站点', '上海市', '上海市', '浦东新区', '上海市浦东新区陆家嘴环路1000号', '021-66666666', '李经理', 1, NOW(), NOW()),
('广州天河站点', '广东省', '广州市', '天河区', '广州市天河区珠江新城花城大道88号', '020-77777777', '陈经理', 1, NOW(), NOW());

-- 插入测试仓库
INSERT INTO `warehouses` (`name`, `location`, `capacity`, `used`, `manager`, `phone`, `status`, `created_at`, `updated_at`)
VALUES 
('北京中心仓库', '北京市朝阳区', 10000, 0, '仓管A', '13800138002', 1, NOW(), NOW()),
('上海中心仓库', '上海市浦东新区', 8000, 0, '仓管B', '13800138003', 1, NOW(), NOW());
