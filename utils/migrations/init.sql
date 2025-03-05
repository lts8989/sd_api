create database if not exists sd_test;
use sd_test;
create table execution_tasks
(
    id          int unsigned auto_increment
        primary key,
    prompt_id   varchar(100)                       not null,
    template_id int unsigned                       not null comment '模版id',
    parameters  json                               not null comment '参数对，使用 JSON 格式存储任务的参数',
    status      tinyint unsigned                   not null comment '任务状态。1、待处理；2、进行中；3、已完成；4、失败',
    result      varchar(500)                       not null comment '执行结果，存储任务执行后的结果信息',
    created_at  datetime default CURRENT_TIMESTAMP not null comment '记录创建时间，默认为当前时间',
    updated_at  datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '记录更新时间，自动更新为当前时间'
);


create table result_images
(
    id         int unsigned auto_increment
        primary key,
    filename   varchar(200)                       not null comment '文件名',
    subfolder  varchar(100)                       not null comment '子文件夹',
    type       varchar(200)                       not null comment '类型',
    prompt_id  varchar(100)                       not null,
    created_at datetime default CURRENT_TIMESTAMP not null
);

create table templates
(
    id            int unsigned auto_increment
        primary key,
    template_name varchar(100)                       not null comment '模板的名称',
    file_name     varchar(100)                       not null comment '与模板相关联的文件名',
    parameters    json                               not null comment '存储与模板相关的参数值',
    created_at    datetime default CURRENT_TIMESTAMP not null
);

INSERT INTO sd_test.templates (id, template_name, file_name, parameters) VALUES (1, 'default_temp', 'default_temp.txt', '{"temp_id":1,"params":{"seed":"4259488860108327","steps":"2","batch_size":"1","prompt_msg":"(((shiny skin:1.3))), (shiny oiled skin), (dark skin:1.1), (detailed face), grin, (chubby:1), 1girl, beautiful woman, young woman, 20yo, (tight sheer white T-shirt), (tight black leggings), (high heels), Very fine and delicate eyes, shiny lips, (huge breasts), breast focus, BREAK ((whole body)), beautiful hip line, ((cleavage)), blonde hair, looking at viewer, BREAK (best quality), 8K, (realistic), ultra high res, extreme detailed, masterpiece, cinematic lighting"}}');

