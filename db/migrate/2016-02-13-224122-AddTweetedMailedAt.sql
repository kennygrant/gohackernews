/* SQL migration AddTweetedMailedAt */
alter table stories add column tweeted_at timestamp;
alter table stories add column newsletter_at timestamp;