import * as React from "react";
import {IMediaResponse, IMediaVideoResponse} from "../pages/media/media-list-page";

interface IProps {
    media: IMediaResponse;
    style?: React.CSSProperties
}

export class Video extends React.Component<IProps, {}> {

    public render(): JSX.Element | null {
        const mp4: string | null = this.getMp4Url();
        if (!mp4) {
            return null;
        }
        return <video controls={false} autoPlay={true} loop={true} style={this.props.style}>
            <source src={mp4} type="video/mp4"/>
        </video>
    }

    private getMp4Url(): string | null {
        const videos: IMediaVideoResponse[] = this.props.media.videos
            .filter((video: IMediaVideoResponse) => video.type === "mp4");
        if (videos.length > 0) {
            return videos[0].url;
        }
        return null;
    }

}