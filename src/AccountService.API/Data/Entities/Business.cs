namespace AccountService.API.Data.Entities;

public class Business : BaseEntity
{
    public Guid UserId { get; set; }
    public required string Name { get; set; }
    public bool Status { get; set; } // TODO: change to isVerified later
    public string? IdentityNumber { get; set; }

    public List<Store> Stores { get; set; } = new();
    public List<Staff> Staffs { get; set; } = new();
}