using AccountService.API.Data;
using AccountService.API.Data.Repositories;
using Microsoft.EntityFrameworkCore;
using MassTransit;
using AccountService.API.Messaging.Consumer;
using Serilog;
using AccountService.API.Endpoints;
using AccountService.API.Services;
using AccountService.API.Handlers;
var builder = WebApplication.CreateBuilder(args);

builder.Host.UseSerilog((context, config) => config.ReadFrom.Configuration(context.Configuration));

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

        var logger = context.GetRequiredService<ILogger<Program>>();
        logger.LogInformation("Connected to RabbitMQ at {Host}", builder.Configuration.GetConnectionString("RabbitMQ"));
    });
});

builder.Services.AddScoped<IBusinessRepository, BusinessRepository>();
builder.Services.AddScoped<IStoreRepository, StoreRepository>();
builder.Services.AddScoped<IBusinessService, BusinessService>();
builder.Services.AddScoped<IStoreService, StoreService>();
builder.Services.AddExceptionHandler<CustomExceptionHandler>();

var app = builder.Build();

using (var scope = app.Services.CreateScope())
{
    var db = scope.ServiceProvider.GetRequiredService<AppDbContext>();
    db.Database.Migrate();
}

app.UseExceptionHandler(options => { });
app.UseSerilogRequestLogging();
app.MapStoreEndpoints();
app.MapBusinessEndpoints();
app.MapGet("/ping", () => Results.Ok("pong"));

app.Run(
    url: "http://+:8000"
);