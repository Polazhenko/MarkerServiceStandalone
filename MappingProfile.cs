using AutoMapper;
using MarkerServiceStandalone.Models;

namespace MarkerServiceStandalone;

public class MappingProfile : Profile
{
    public MappingProfile()
    {
        CreateMap<Marker, MarkerResponse>();
        CreateMap<CreateMarkerRequest, Marker>();
    }
}
