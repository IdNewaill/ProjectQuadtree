import os
from tkinter import Tk, filedialog
from pydub import AudioSegment

def convert_mp3_to_wav(input_folder):
    # Parcours tous les fichiers dans le dossier et les sous-dossiers
    for root, dirs, files in os.walk(input_folder):
        for filename in files:
            if filename.lower().endswith(".mp3"):
                mp3_path = os.path.join(root, filename)
                wav_path = os.path.join(root, f"{os.path.splitext(filename)[0]}.wav")
                
                print(f"Conversion de {filename} en WAV...")
                
                try:
                    # Chargement du fichier MP3 et conversion en WAV
                    audio = AudioSegment.from_mp3(mp3_path)
                    audio.export(wav_path, format="wav")
                    print(f"{filename} converti avec succès en {os.path.basename(wav_path)}")
                    
                    # Suppression du fichier MP3 après conversion
                    os.remove(mp3_path)
                    print(f"Le fichier MP3 {filename} a été supprimé.")
                
                except Exception as e:
                    print(f"Erreur lors de la conversion de {filename}: {e}")

def select_folder():
    # Crée une fenêtre Tkinter sans interface utilisateur
    root = Tk()
    root.withdraw()  # Ne pas afficher la fenêtre principale
    # Ouvre la boîte de dialogue pour choisir un dossier
    folder_selected = filedialog.askdirectory(title="Sélectionner un dossier contenant des fichiers MP3")
    return folder_selected

if __name__ == "__main__":
    folder = select_folder()
    if folder:
        convert_mp3_to_wav(folder)
        print("Conversion terminée.")
    else:
        print("Aucun dossier sélectionné.")
