namespace AccountService.API.Data.Entities;

public class Staff: BaseEntity
{
    public Guid UserId { get; set; } // from auth
    public Guid BusinessId { get; set; }

    public Business? Business { get; set; }
    public List<StaffPermission> Permissions { get; set; } = new();
    public List<StaffStore> StaffStores { get; set; } = new();
}


public class StaffPermission
{
    public Guid Id { get; set; }
    public Guid StaffId { get; set; }
    public string ActionKey { get; set; } = default!;

    public Staff? Staff { get; set; }
}

public class StaffStore
{
    public Guid StaffId { get; set; }
    public Guid StoreId { get; set; }

    public Staff? Staff { get; set; }
    public Store? Store { get; set; }
}
