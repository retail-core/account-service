using AccountService.Core.Interfaces;
using AccountService.Core.Services;
using AccountService.Data;
using AccountService.Data.Repositories;
using Microsoft.EntityFrameworkCore;
using MassTransit;
using AccountService.Messaging.Consumer;
using AccountService.Core.Entities;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddDbContext<AppDbContext>(options =>
    options.UseNpgsql(builder.Configuration.GetConnectionString("Postgres")));

builder.Services.AddMassTransit(x =>
{
    x.AddConsumer<OwnerCreatedConsumer>();

    x.UsingRabbitMq((context, cfg) =>
    {
        cfg.Host(builder.Configuration.GetConnectionString("RabbitMQ"));
        cfg.ReceiveEndpoint("account-service", e =>
        {
            e.ConfigureConsumeTopology = false;
            e.UseRawJsonSerializer(isDefault: true);
            e.ConfigureConsumer<OwnerCreatedConsumer>(context);

            e.Bind("domain.events", s =>
            {
                s.RoutingKey = "business_owner.created";
                s.ExchangeType = "topic";
            });
        });
    });
});

builder.Services.AddScoped<IBaseRepository<Business>, BusinessRepository>();
builder.Services.AddScoped<IBusinessService, BusinessService>();
builder.Services.AddScoped<IBaseRepository<Store>, StoreRepository>();
builder.Services.AddScoped<IStoreService, StoreService>();

var app = builder.Build();

app.MapGet("/status", () => Results.Ok("Ok"));

app.Run();






// if (app.Environment.IsDevelopment())
// {
//     app.UseSwagger();
//     app.UseSwaggerUI();
// }
//builder.Services.AddSwaggerGen();
// builder.Services.AddEndpointsApiExplorer();