import type { Service, Tag } from "@/types/service";

export interface Config {
    Title: string;
    Services: Service[];
    Tags: Tag[];
}
