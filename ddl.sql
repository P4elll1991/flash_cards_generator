create schema flash_cards;

create table flash_cards.tasks (
    id serial,
    create_date timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    "status" varchar DEFAULT 'NEW',
    details text
);

create table flash_cards.task_params (
    task_id integer,
    native_lang varchar,
    learning_lang varchar,
    words integer,
    topic varchar,
    "level" varchar
);

create table flash_cards.cards (
    task_id integer,
    word varchar,
    native_lang varchar,
    learning_lang varchar,
    topic varchar,
    "level" varchar,
    pronunciation text,
    phonetic_respelling text,
    "definition" text,
    "translation" text,
    example text,
    example_translation text
);
