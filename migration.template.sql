create table if not exists {{.Name}} (
  id         integer primary key,
  created_at timestamp not null default current_timestamp,
  updated_at timestamp not null default current_timestamp,
  deleted_at timestamp,
  data       blob,
  type       text not null,
  owner_id   integer,
  ownership  text,

  foreign key (owner_id) references {{.Name}}(id)
);

create index if not exists ix_{{.Name}}_created_at
             on {{.Name}} (created_at);
create index if not exists ix_{{.Name}}_updated_at
             on {{.Name}} (updated_at);
create index if not exists ix_{{.Name}}_deleted_at
             on {{.Name}} (deleted_at)
             where deleted_at is not null;
create index if not exists ix_{{.Name}}_type_id
             on {{.Name}} (type);
create index if not exists ix_{{.Name}}_owner_id
             on {{.Name}} (owner_id);
create index if not exists ix_{{.Name}}_ownership
             on {{.Name}} (ownership)
             where owner_id is not null;

create table if not exists {{.Name}}_values (
  id          integer primary key,
  field       text not null,
  value       blob not null,
  entity_type text not null,
  entity_id   integer not null,
  is_indexed  boolean not null default FALSE,
  is_unique   boolean not null default FALSE,

  foreign key (entity_id) references {{.Name}}(id)
);

create index if not exists ix_{{.Name}}_values_field
             on entity_values (field);
create index if not exists ix_{{.Name}}_values_entity_id
             on entity_values (entity_id);
create index if not exists ix_{{.Name}}_values_field_value
             on entity_values (field, value)
             where is_indexed;
create unique index if not exists uk_{{.Name}}_values_field_value_entity_type
             on entity_values (entity_type, field, value)
             where is_unique;

create table if not exists {{.Name}}_relationships (
  id          integer primary key,
  name        text not null,
  owner_type  text not null,
  owner_id    integer not null,
  entity_type text not null,
  entity_id   integer not null,
  data        blob,

  foreign key (owner_id) references {{.Name}}(id),
  foreign key (entity_id) references {{.Name}}(id)
);

create index if not exists ix_{{.Name}}_relationships_name
             on relationships (name);
create index if not exists ix_{{.Name}}_relationships_owner_type
             on relationships (owner_type);
create index if not exists ix_{{.Name}}_relationships_owner_id
             on relationships (owner_id);
create index if not exists ix_{{.Name}}_relationships_entity_type
             on relationships (entity_type);
create index if not exists ix_{{.Name}}_relationships_entity_id
             on relationships (entity_id);
