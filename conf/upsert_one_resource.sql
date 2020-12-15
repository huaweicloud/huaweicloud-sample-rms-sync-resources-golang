INSERT INTO resources (provider, type, resourceId, name, regionId, projectId, projectName, epId, epName, checksum, created, updated, provisioningState, tags, properties, queryAt, state)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
name = IF(checksum <> VALUES(checksum) AND queryAt < VALUES(queryAt), VALUES(name), name),
regionId = IF(checksum <> VALUES(checksum) AND queryAt < VALUES(queryAt), VALUES(regionId), regionId),
projectId = IF(checksum <> VALUES(checksum) AND queryAt < VALUES(queryAt), VALUES(projectId), projectId),
projectName = IF(checksum <> VALUES(checksum) AND queryAt < VALUES(queryAt), VALUES(projectName), projectName),
epId = IF(checksum <> VALUES(checksum) AND queryAt < VALUES(queryAt), VALUES(epId), epId),
epName = IF(checksum <> VALUES(checksum) AND queryAt < VALUES(queryAt), VALUES(epName), epName),
checksum = IF(checksum <> VALUES(checksum) AND queryAt < VALUES(queryAt), VALUES(checksum), checksum),
created = IF(checksum <> VALUES(checksum) AND queryAt < VALUES(queryAt), VALUES(created), created),
updated = IF(checksum <> VALUES(checksum) AND queryAt < VALUES(queryAt), VALUES(updated), updated),
provisioningState = IF(checksum <> VALUES(checksum) AND queryAt < VALUES(queryAt), VALUES(provisioningState), provisioningState),
tags = IF(checksum <> VALUES(checksum) AND queryAt < VALUES(queryAt), VALUES(tags), tags),
properties = IF(checksum <> VALUES(checksum) AND queryAt < VALUES(queryAt), VALUES(properties), properties),
queryAt = IF(checksum <> VALUES(checksum) AND queryAt < VALUES(queryAt), VALUES(queryAt), queryAt),
state = IF(checksum <> VALUES(checksum) AND queryAt < VALUES(queryAt), VALUES(state), state)