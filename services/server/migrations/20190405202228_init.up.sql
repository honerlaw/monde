CREATE TABLE IF NOT EXISTS `user` (
  `id` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `username` varchar(255) NOT NULL,
  `hash` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
  KEY `idx_user_deleted_at` (`deleted_at`),
  KEY `idx_user_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `media` (
  `id` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `user_id` varchar(255) DEFAULT NULL,
  `job_id` varchar(255) DEFAULT NULL,
  `title` varchar(255) DEFAULT NULL,
  `description` text DEFAULT NULL,
  `published` bool DEFAULT FALSE,
  `published_data` timestamp NULL default NULL,
  PRIMARY KEY (`id`),
  KEY `idx_media_deleted_at` (`deleted_at`),
  KEY `media_user_id_user_id_foreign` (`user_id`),
  CONSTRAINT `media_user_id_user_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `track` (
  `id` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `media_id` varchar(255) DEFAULT NULL,
  `type` varchar(255) DEFAULT NULL,
  `duration` double DEFAULT NULL,
  `width` bigint(20) DEFAULT NULL,
  `height` bigint(20) DEFAULT NULL,
  `format` varchar(255) DEFAULT NULL,
  `encoded_date` varchar(255) DEFAULT NULL,
  `video_count` varchar(255) DEFAULT NULL,
  `data_size` bigint(20) DEFAULT NULL,
  `file_size` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_track_deleted_at` (`deleted_at`),
  KEY `track_media_id_media_id_foreign` (`media_id`),
  CONSTRAINT `track_media_id_media_id_foreign` FOREIGN KEY (`media_id`) REFERENCES `media` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `hashtag` (
  `id` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `tag` tinytext NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_media_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `media_hashtag` (
  `media_id` varchar(255) NOT NULL,
  `hashtag_id` varchar(255) NOT NULL,
  PRIMARY KEY (`media_id`, `hashtag_id`),
  KEY `media_hashtag_media_id_foreign` (`media_id`),
  KEY `media_hashtag_hashtag_id_foreign` (`hashtag_id`),
  CONSTRAINT `media_hashtag_media_id_foreign` FOREIGN KEY (`media_id`) REFERENCES `media` (`id`) ON DELETE CASCADE,
  CONSTRAINT `media_hashtag_hashtag_id_foreign` FOREIGN KEY (`hashtag_id`) REFERENCES `hashtag` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;