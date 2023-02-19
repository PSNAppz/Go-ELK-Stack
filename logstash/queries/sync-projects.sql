SELECT pl.id,
       pl.operation,
       pl.project_id,
       p.id,
       p.name,
       p.slug,
       p.description,
       p.created_at,
       u.id as user_id,
       u.name as user_name,
       array_agg(h.name) as hashtags
FROM project_logs pl
        LEFT JOIN projects p ON p.id = pl.project_id
        LEFT JOIN user_projects up ON p.id = up.project_id
        LEFT JOIN users u ON up.user_id = u.id
        LEFT JOIN project_hashtags ph ON ph.project_id = p.id
        LEFT JOIN hashtags h ON h.id = ph.hashtag_id
WHERE pl.id > :sql_last_value 
GROUP BY pl.id, p.id, u.id 
ORDER BY pl.id;
