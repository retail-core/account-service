namespace AccountService.API.Messaging.Contracts;

public record OwnerCreatedEvent(Guid OwnerId, string OwnerName);