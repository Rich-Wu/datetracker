interface DateData {
    id: string;
    age: number;
    cost?: number | null;
    date: string;
    createdAt: string;
    ethnicity: string[];
    firstName: string;
    lastName?: string | null;
    ownerId: string;
    places: Place[];
    result: string;
    split: boolean;
}

interface Place {
    place: string;
    typeOfPlace: string;
    cost: number;
}