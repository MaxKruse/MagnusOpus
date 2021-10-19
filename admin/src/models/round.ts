class Round {
    id: number;
    name: string;
    description: string;
    start_time: string;
    end_time: string;
    download_path: string;

    constructor(id: number, name: string, description: string, start_time: string, end_time: string, download_path: string) {
        this.id = id;
        this.name = name;
        this.description = description;
        this.start_time = start_time;
        this.end_time = end_time;
        this.download_path = download_path;
    }
}

export default Round;