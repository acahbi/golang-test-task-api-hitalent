-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.department(
      id            BIGSERIAL PRIMARY KEY -- Если у нас будет очень много id, лучше использовать BIGSERIAL
    , name          VARCHAR(200) NOT NULL
    , parent_id     BIGINT
    , created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
    , updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

COMMENT ON TABLE public.department              IS 'Таблица подразделений';
COMMENT ON COLUMN public.department.id          IS 'Уникальный идентификатор подразделения';
COMMENT ON COLUMN public.department.name        IS 'Имя подразделения';
COMMENT ON COLUMN public.department.parent_id   IS 'Родительский ID';
COMMENT ON COLUMN public.department.created_at  IS 'Время создания записи';
COMMENT ON COLUMN public.department.updated_at  IS 'Время обновления записи';

ALTER TABLE public.department ADD CONSTRAINT department_unique_parent_id_name UNIQUE (parent_id, name);
ALTER TABLE public.department ADD CONSTRAINT department_parent_id FOREIGN KEY (parent_id) REFERENCES public.department(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.department;
-- +goose StatementEnd