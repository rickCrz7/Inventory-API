create table owners (
    id varchar(50) primary key,
    first_name varchar(50) not null,
    last_name varchar(50) not null,
    campus_id varchar(50),
    email varchar(50) not null
);

create table types (
    id varchar(50) primary key,
    name varchar(50) not null,
    description text
);

create table type_properties (
    id varchar(50) primary key,
    type_id varchar(50) not null,
    name varchar(50) not null,
    data_type varchar(50) not null,
    required boolean not null,
    foreign key (type_id) references types(id) on delete cascade
);

create index idx_type_properties_type_id on type_properties(type_id);

create table devices (
    id varchar(50) primary key,
    serial_number varchar(50),
    name varchar(50) not null,
    type_id varchar(50) not null,
    owner_id varchar(50) not null,
    purchase_date date,
    status varchar(50) not null,
    foreign key (type_id) references types(id),
    foreign key (owner_id) references owners(id) on delete cascade
);

create index idx_devices_type_id on devices(type_id);
create index idx_devices_owner_id on devices(owner_id);

create table device_properties (
    id varchar(50) primary key,
    device_id varchar(50) not null,
    type_property_id varchar(50) not null,
    value text not null,
    foreign key (device_id) references devices(id) on delete cascade,
    foreign key (type_property_id) references type_properties(id) on delete cascade
);

create index idx_device_properties_device_id on device_properties(device_id);
create index idx_device_properties_type_property_id on device_properties(type_property_id);

create table device_photos (
    id varchar(50) primary key,
    device_id varchar(50) not null,
    photo text not null,
    created_at timestamp not null,
    foreign key (device_id) references devices(id) on delete cascade
);

create index idx_device_photos_device_id on device_photos(device_id);

create table device_logs (
    id varchar(50) primary key,
    device_id varchar(50) not null,
    log_type varchar(50) not null,
    note text,
    created_at timestamp not null,
    created_by varchar(50) not null,
    foreign key (device_id) references devices(id) on delete cascade
);

create index idx_device_logs_device_id on device_logs(device_id);

-- create table device_assignments (
--     id varchar(50) primary key,
--     device_id varchar(50) not null,
--     owner_id varchar(50) not null,
--     assigned_at timestamp not null,
--     returned_at timestamp,
--     foreign key (device_id) references devices(id) on delete cascade,
--     foreign key (owner_id) references owners(id) on delete cascade
-- );

-- create index idx_device_assignments_device_id on device_assignments(device_id);
-- create index idx_device_assignments_owner_id on device_assignments(owner_id);

-- create table device_assignment_history (
--     id varchar(50) primary key,
--     device_assignment_id varchar(50) not null,
--     status varchar(50) not null,
--     changed_at timestamp not null,
--     changed_by varchar(50) not null,
--     foreign key (device_assignment_id) references device_assignments(id) on delete cascade
-- );

-- create index idx_device_assignment_history_device_assignment_id on device_assignment_history(device_assignment_id);

-- drop table if exists device_assignment_history;
-- drop table if exists device_assignments;

drop table if exists device_logs;
drop table if exists device_photos;
drop table if exists device_properties;
drop table if exists devices;
drop table if exists type_properties;
drop table if exists types;
drop table if exists owners;