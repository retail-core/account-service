using System;

namespace AccountService.API.Data.Entities;

public class Store : BaseEntity
{
    public Guid BusinessId { get; set; }
    public required string Name { get; set; }
    public string? Address { get; set; }
    public string? City { get; set; }
    public string? State { get; set; }
    public string? Country { get; set; }

    public Business? Business { get; set; }
    public List<StaffStore> StaffStores { get; set; } = new();
}
