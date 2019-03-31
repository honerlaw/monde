import * as React from "react";
import {IMediaResponse, IMediaVideoResponse} from "../pages/upload-list-page";

interface IProps {
    media: IMediaResponse;
    isLink: boolean;

}

export class ThumbVideo extends React.Component<IProps, {}> {

    public render(): JSX.Element | null {
        const mp4: string | null = this.getMp4Url();
        if (this.props.media.thumbnails.length === 0 || !mp4) {
            return null;
        }

        if (this.props.isLink) {
            return <a href={`/media/view/${this.props.media.id}`}>
                {this.renderMainView(mp4)}
            </a>
        }

        return this.renderMainView(mp4);
    }

    private renderMainView(mp4: string): JSX.Element {
        return <div className={"media-thumbnail-video"}>
            <div className={"media-thumbnail-video-thumbnail"}>
                <div className="img" style={this.getImageStyle()}/>
            </div>
            <video controls={false} autoPlay={true} loop={true}>
                <source src={mp4} type="video/mp4"/>
            </video>
        </div>;
    }

    private getImageStyle(): React.CSSProperties {
        return {
            backgroundImage: `url("${this.props.media.thumbnails[0]}")`,
            height: `${this.props.media.videos[0].height}px`,
            width: `${this.props.media.videos[0].width}px`
        }
    }

    private getMp4Url(): string | null {
        const videos: IMediaVideoResponse[] = this.props.media.videos.filter((video) => video.type === "mp4");
        if (videos.length > 0) {
            return videos[0].url;
        }
        return null;
    }

}