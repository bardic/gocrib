git log --pretty=format:user:%aN%n%ct \
    --reverse \
    --raw \
    --encoding=UTF-8 \
    --no-renames \
    --since="2024-10-15" \
    > gitcommits.log

gource \
    gitcommits.log \
    -1280x720 \
    --output-ppm-stream - \
    --output-framerate 30 \
    --seconds-per-day 2 \
    --stop-at-end \
    | ffmpeg -y -r 30 -f image2pipe -vcodec ppm -i - -b 65536K movie.mp4