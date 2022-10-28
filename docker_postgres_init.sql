CREATE TABLE
    IF NOT EXISTS public.users (
        id integer NOT NULL GENERATED ALWAYS AS IDENTITY,
        name character varying(50) NOT NULL,
        PRIMARY KEY (id)
    );

CREATE TABLE
    IF NOT EXISTS public.main_account (
        id integer NOT NULL GENERATED ALWAYS AS IDENTITY,
        user_id integer NOT NULL,
        bill_amount integer NOT NULL,
        PRIMARY KEY (id),
        UNIQUE (user_id)
    );

CREATE TABLE
    IF NOT EXISTS public.reserve_account (
        id integer NOT NULL GENERATED ALWAYS AS IDENTITY,
        user_id integer NOT NULL,
        bill_amount integer NOT NULL,
        PRIMARY KEY (id),
        UNIQUE (user_id)
    );

CREATE TABLE
    IF NOT EXISTS public.transactions (
        id integer NOT NULL GENERATED ALWAYS AS IDENTITY,
        user_id integer NOT NULL,
        service_id integer,
        order_id integer,
        start_date timestamp
        with
            time zone NOT NULL,
            end_date timestamp
        with
            time zone,
            start_bill_amount integer NOT NULL,
            end_bill_amount integer NOT NULL,
            PRIMARY KEY (id)
    );

ALTER TABLE
    IF EXISTS public.main_account
ADD
    FOREIGN KEY (user_id) REFERENCES public.users (id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION NOT VALID;

ALTER TABLE
    IF EXISTS public.reserve_account
ADD
    FOREIGN KEY (user_id) REFERENCES public.users (id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION NOT VALID;

ALTER TABLE
    IF EXISTS public.transactions
ADD
    FOREIGN KEY (user_id) REFERENCES public.users (id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION NOT VALID;