SELECT pl.id,
       pl.operation,
       pl.project_id,
       p.id,
       p.name,
       p.slug,
       p.description
FROM project_logs pl
         LEFT JOIN projects p
                   ON p.id = pl.project_id
WHERE pl.id > :sql_last_value ORDER BY pl.id;
