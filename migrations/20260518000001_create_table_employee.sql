-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.employee(
      id            BIGSERIAL PRIMARY KEY
    , department_id BIGINT
    , full_name     VARCHAR(200) NOT NULL
    , position      VARCHAR(200) NOT NULL
    , hired_at      TIMESTAMP
    , created_at    TIMESTAMP
    , updated_at    TIMESTAMP
);

COMMENT ON TABLE public.employee                IS 'Таблица подразделений';
COMMENT ON COLUMN public.employee.id            IS 'Уникальный идентификатор подразделения';
COMMENT ON COLUMN public.employee.full_name     IS 'Полное имя подразделения';
COMMENT ON COLUMN public.employee.position      IS 'Позиция/Должность';
COMMENT ON COLUMN public.employee.hired_at      IS 'Дата приёма на работу';
COMMENT ON COLUMN public.employee.created_at    IS 'Время создания записи';
COMMENT ON COLUMN public.employee.updated_at    IS 'Время обновления записи';

CREATE INDEX employee_department_id_ind ON public.employee(department_id);
ALTER TABLE public.employee ADD CONSTRAINT employee_department_id_fk FOREIGN KEY (department_id) REFERENCES public.department(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.employee;
-- +goose StatementEnd