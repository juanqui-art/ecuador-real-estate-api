'use client';

import { 
  Dialog, 
  DialogContent, 
  DialogHeader, 
  DialogTitle,
  DialogDescription
} from '@/components/ui/dialog';
import { EditPropertyForm } from '@/components/forms/edit-property-form';
import { ScrollArea } from '@/components/ui/scroll-area';

interface Property {
  id: string;
  title: string;
  description: string;
  price: number;
  type: string;
  status: string;
  province: string;
  city: string;
  address: string | null;
  bedrooms: number;
  bathrooms: number;
  area_m2: number;
  parking_spaces: number;
  year_built?: number | null;
  has_garden: boolean;
  has_pool: boolean;
  has_elevator: boolean;
  has_balcony: boolean;
  has_terrace: boolean;
  has_garage: boolean;
  is_furnished: boolean;
  allows_pets: boolean;
  contact_phone: string | null;
  contact_email: string | null;
  notes?: string | null;
  created_at: string;
  updated_at: string;
  is_featured?: boolean;
  images?: string[];
  main_image?: string | null;
}

interface EditPropertyDialogProps {
  property: Property | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: () => void;
}

export function EditPropertyDialog({ 
  property, 
  open, 
  onOpenChange, 
  onSuccess 
}: EditPropertyDialogProps) {
  if (!property) return null;

  const handleSuccess = () => {
    onSuccess?.();
    onOpenChange(false);
  };

  const handleCancel = () => {
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-4xl max-h-[90vh] p-0">
        <DialogHeader className="p-6 pb-0">
          <DialogTitle className="text-xl">
            Editar Propiedad
          </DialogTitle>
          <DialogDescription>
            Modifica los datos de la propiedad "{property.title}"
          </DialogDescription>
        </DialogHeader>
        
        <ScrollArea className="max-h-[calc(90vh-120px)] px-6 pb-6">
          <EditPropertyForm 
            property={property}
            onSuccess={handleSuccess}
            onCancel={handleCancel}
          />
        </ScrollArea>
      </DialogContent>
    </Dialog>
  );
}