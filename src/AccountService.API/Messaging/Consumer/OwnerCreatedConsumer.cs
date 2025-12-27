using AccountService.API.Services;
using AccountService.API.Messaging.Contracts;
using MassTransit;

namespace AccountService.API.Messaging.Consumer;

public class OwnerCreatedConsumer(IBusinessService _bizService, IStoreService _storeService, ILogger<OwnerCreatedConsumer> logger) : IConsumer<OwnerCreatedEvent>
{
    public async Task Consume(ConsumeContext<OwnerCreatedEvent> context)
    {
        var message = context.Message;
        logger.LogInformation("Received OwnerCreatedEvent for OwnerId: {OwnerId}, OwnerName: {OwnerName}", message.OwnerId, message.OwnerName);
        var biz = await _bizService.CreateBusinessAsync(message.OwnerName + " Business", message.OwnerId);
        logger.LogInformation("Business created for OwnerId: {OwnerId}", message.OwnerId);
        await _storeService.CreateStoreAsync(message.OwnerName + " Main Store", biz!.Id);
    }
}