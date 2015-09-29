/* SQL migration AddCommentLevels */
alter table comments add column rank integer;
alter table comments add column level integer;
alter table comments add column dotted_ids text; 