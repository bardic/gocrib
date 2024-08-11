INSERT INTO "items" ("id", "name", "cost", "weight", "unit", "barcode", "storeid", "storename", "storeneighborhood", "tags", "created_at", "updated_at") VALUES
(1,	'test1',	1,	1,	'1',	'1',	'store1',	'store1',	'here',	'{1}',	'2024-05-03 11:25:47.576503+00',	'2024-05-03 11:25:47.576503+00'),
(2,	'test2',	2,	2,	'1',	'2',	'store1',	'store1',	'here',	'{1,2}',	'2024-05-03 11:26:11.294413+00',	'2024-05-03 11:26:11.294413+00'),
(3,	'test3',	2,	2,	'1',	'3',	'store1',	'store1',	'here',	'{1,2}',	'2024-05-03 11:26:32.262067+00',	'2024-05-03 11:26:32.262067+00'),
(4,	'test4',	3,	2,	'1',	'4',	'store1',	'store1',	'here',	'{2,3}',	'2024-05-03 11:26:52.160195+00',	'2024-05-03 11:26:52.160195+00');

INSERT INTO "tags" ("id", "name") VALUES
(1,	'fruit'),
(2,	'meat'),
(3,	'snack');