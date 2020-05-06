FROM sighery/archbuilder:latest

COPY --chown=builder ./ package/
WORKDIR package
RUN makepkg -s --needed --noconfirm --noprogressbar

CMD ["/bin/bash"]
