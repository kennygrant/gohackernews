/* SQL migration AddUserNames */

alter table comments add column user_name text;
alter table stories add column user_name text;
alter table comments add column story_name text;

alter table stories add column comment_count integer;

alter table comments drop column created_by;
alter table stories drop column created_by;

alter table comments add column user_id integer;
alter table stories add column user_id integer;