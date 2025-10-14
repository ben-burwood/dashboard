
export interface Icon {
  Name: string;
  Color?: string;
};

export interface Tag {
	Name : string;
    Color: string;
}

export interface Service {
  Link: string;
  Icon: Icon;
  Title: string;
  Description?: string;
  Group?: string;
  Tags?: string[];
};
